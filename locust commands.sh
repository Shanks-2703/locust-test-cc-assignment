git clone 


cd terraform-google-three-tier-web-app

gcloud run services describe three-tier-app-api --region asia-south1 --project trusty-stacker-453107-i1   #replace tiered-web-app-api with the cloud run instance of backend API

gcloud run services describe three-tier-app-fe --region asia-south1 --project trusty-stacker-453107-i1    #replace tiered-web-app-fe with the cloud run instance of Frontend

gcloud run services update three-tier-app-fe --region asia-south1 --project trusty-stacker-453107-i1 --set-env-vars "ENDPOINT=https://three-tier-app-api-1049385999004.asia-south1.run.app"   #replace ENDPOINT URL with the cloud run instance of backend URL

gcloud sql instances patch three-tier-app-db-4097 --assign-ip --project three-tier-web-app-457409

# grab your Cloud Shell external IP
export MY_IP=$(curl -s https://ifconfig.me)
gcloud sql instances patch three-tier-app-db-4097 --authorized-networks="${MY_IP}/32" --project three-tier-web-app-457409

wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O cloud_sql_proxy

chmod +x cloud_sql_proxy

pkill cloud_sql_proxy 

./cloud_sql_proxy \
  -instances="three-tier-web-app-457409:us-central1:three-tier-app-db-4097=tcp:5432" \
  -ip_address_types=PUBLIC &

gcloud sql users set-password postgres --instance=three-tier-app-db-4097 --password=Newpassword --project=three-tier-web-app-457409

psql "host=127.0.0.1 port=5432 dbname=todo user=postgres password=Newpassword sslmode=disable"

CREATE TABLE loadtest_table (
  id   SERIAL PRIMARY KEY,
  stub TEXT
);

-- 2. Seed a few rows
INSERT INTO loadtest_table (stub) VALUES
  ('foo'), ('bar'), ('baz');

GRANT USAGE   ON SCHEMA public TO postgres;
GRANT SELECT  ON TABLE  loadtest_table TO postgres;

locust --headless --users 200 --spawn-rate 5 --run-time 3m -f locustfile.py --only-summary
