import http from "k6/http";

const BASE_URL = "http://nginx";

export default function () {
    http.get(`${BASE_URL}`);
}
