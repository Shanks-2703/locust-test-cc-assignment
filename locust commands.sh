git clone 


cd terraform-google-three-tier-web-app

gcloud run services describe tiered-web-app-api --reigion asia-south1    #replace tiered-web-app-api with the cloud run instance of backend API

gcloud run services describe tiered-web-app-fe --reigion asia-south1     #replace tiered-web-app-fe with the cloud run instance of Frontend

gcloud run services update tiered-web-app-fe --region asia-south1 --set-env-vars "ENDPOINT=https://tiered-web-app-api-880158088034.asia-south1.run.app"   #replace ENDPOINT URL with the cloud run instance of backend URL

gcloud sql instances patch three-tier-app-db-8eef \
  --assign-ip

# grab your Cloud Shell external IP
export MY_IP=$(curl -s https://ifconfig.me)
gcloud sql instances patch  tiered-web-app-db-044f \
  --authorized-networks="${MY_IP}/32"

pkill cloud_sql_proxy 

./cloud_sql_proxy \
  -instances="sonorous-saga-458113-n7:asia-south1: tiered-web-app-db-044f=tcp:5432" \
  -ip_address_types=PUBLIC &

gcloud sql users set-password postgres --instance=tiered-web-app-db-044f --password=Newpassword --project=sonorous-saga-458113-n7

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

locust   --headless   --users  200   --spawn-rate 5   --run-time 3m   -f locustfile.py   --only-summary
