import http from "k6/http";
import { check } from "k6";

export let options = { vus: 50, iterations: 200 };

export default function () {
    const payload = JSON.stringify({
        pull_request_id: "pr-1001",
        pull_request_name: "Add login",
        author_id: "u1"
    });

    const res = http.post("http://avito-test-assignment:8080/pullRequest/create", payload, {
        headers: { "Content-Type": "application/json" },
    });

    console.log("CREATE:", res.status, res.body);

    check(res, {
        "201 created": r => r.status === 201 || r.status === 409,
    });
}
