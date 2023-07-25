#!/usr/bin/env python3

import os
import json
import sys
from http.server import BaseHTTPRequestHandler, HTTPServer

class MockServerHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self._handle_request("GET")

    def do_POST(self):
        self._handle_request("POST")

    def do_PUT(self):
        self._handle_request("PUT")

    def do_DELETE(self):
        self._handle_request("DELETE")

    def _handle_request(self, method):
        base_path = self.server.base_path
        path = self.path.strip("/")
        if path == "":
            self.send_response(200)
            self.send_header("Content-type", "text/plain")
            self.end_headers()
            self.wfile.write(b"Welcome! I am ready for you.")
            return

        file_path = os.path.join(base_path, f"{path}.{method.lower()}.json").strip("/")

        print("Looking for", file_path)

        if not os.path.isfile(file_path):
            file_path = os.path.join(base_path, f"{path}.json").strip("/")
            print("File not found, trying ", file_path)

        if os.path.isfile(file_path):
            self.send_response(200)
            self.send_header("Content-type", "application/json")
            self.end_headers()
            with open(file_path, "r") as file:
                data = file.read()
                self.wfile.write(data.encode())
        else:
            self.send_response(404)
            self.send_header("Content-type", "text/plain")
            self.end_headers()
            self.wfile.write(b"Sorry, I can't find that file.")

def run_mock_server(base_path, port=8080):
    server_address = ("localhost", port)
    httpd = HTTPServer(server_address, MockServerHandler)
    httpd.base_path = base_path
    print(f"Mock server is running at http://localhost:{port}/")
    httpd.serve_forever()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python mock_server.py [base_path]")
        sys.exit(1)

    base_path = sys.argv[1]
    if not os.path.isdir(base_path):
        print("Error: The specified base path does not exist.")
        sys.exit(1)

    run_mock_server(base_path)

