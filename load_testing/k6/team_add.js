import http from "k6/http";
import { check } from "k6";

export let options = { vus: 50, iterations: 200 };

export default function () {
    const url = "http://avito-test-assignment:8080/team/add";

    const payload = JSON.stringify({
        team_name: "backend",
        members: [
            { user_id: "u1", username: "Alice", is_active: true },
            { user_id: "u2", username: "Bob", is_active: true },
            { user_id: "u3", username: "Charlie", is_active: true },
        ]
    });

    const res = http.post(url, payload, {
        headers: { "Content-Type": "application/json" },
    });

    console.log("STATUS:", res.status, "BODY:", res.body);

    check(res, {
        "201 created": r => r.status === 201 || r.status === 400,
    });
}
