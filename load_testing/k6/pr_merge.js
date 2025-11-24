import http from "k6/http";
import { check } from "k6";

export let options = { vus: 50, iterations: 200 };

export default function () {
    const payload = JSON.stringify({
        pull_request_id: "pr-1001"
    });

    const res = http.post("http://avito-test-assignment:8080/pullRequest/merge", payload, {
        headers: { "Content-Type": "application/json" },
    });

    check(res, {
        "200 merged": r => r.status === 200 || r.status === 404,
    });
}
