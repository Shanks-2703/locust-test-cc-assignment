from locust import HttpUser, User, task, between
import psycopg2
import random
import time

DB_CONFIG = {
    "host": "127.0.0.1",  # Replace with Cloud SQL public IP or proxy
    "port": 5432,
    "dbname": "todo",
    "user": "postgres",
    "password": "YY]3mq(m?oLlyMK("
}

class FrontendUser(HttpUser):
    host = "https://three-tier-app-fe-qlxblrcnua-el.a.run.app"
    wait_time = between(1,3)

    @task(3)
    def index(self):
        self.client.get("/", name="GET /")


class TodoUser(HttpUser):
    wait_time = between(1, 3)
    host = "https://three-tier-app-api-815139404174.asia-south1.run.app"

    @task(4)
    def list_todos(self):
        self.client.get("/api/v1/todo", name="GET /api/v1/todo")

    @task(1)
    def create_todo(self):
        title = f"Load test task #{random.randint(1,10000)}"
        self.client.post(
            "/api/v1/todo",
            data={"title": title},
            name="POST /api/v1/todo",
            timeout=10
        )


class DBLoadUser(User):
    wait_time = between(1, 3)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.conn   = psycopg2.connect(**DB_CONFIG)
        self.cursor = self.conn.cursor()

    @task
    def query_todos(self):
        start = time.time()
        try:
            self.cursor.execute("SELECT * FROM loadtest_table LIMIT 10")
            self.conn.commit()
            rt = int((time.time() - start) * 1000)
            # unified event: exception=None indicates success
            self.environment.events.request.fire(
                request_type="db",
                name="SELECT loadtest_table",
                response_time=rt,
                response_length=0,
                exception=None
            )
        except Exception as e:
            rt = int((time.time() - start) * 1000)
            # exception!=None indicates failure
            self.environment.events.request.fire(
                request_type="db",
                name="SELECT loadtest_table",
                response_time=rt,
                response_length=0,
                exception=e
            )

    def on_stop(self):
        self.cursor.close()
        self.conn.close()