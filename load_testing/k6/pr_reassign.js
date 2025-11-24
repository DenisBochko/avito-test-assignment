import http from "k6/http";
import { check } from "k6";

export let options = { vus: 50, iterations: 200 };

export default function () {
    const payload = JSON.stringify({
        pull_request_id: "pr-1001",
        old_user_id: "u2"
    });

    const res = http.post("http://avito-test-assignment:8080/pullRequest/reassign", payload, {
        headers: { "Content-Type": "application/json" },
    });

    console.log("REASSIGN:", res.status, res.body);

    check(res, {
        "200 OK or 409 domain error": r => [200,404,409].includes(r.status),
    });
}
