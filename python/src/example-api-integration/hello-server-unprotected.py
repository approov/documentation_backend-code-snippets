# @link https://approov.io/docs/v2.0/approov-usage-documentation/#backend-integration-impact
from flask import Flask
app = Flask(__name__)

@app.route("/")
def hello():
   return "Hello World!\n"

if __name__ == "__main__":
   app.run()
