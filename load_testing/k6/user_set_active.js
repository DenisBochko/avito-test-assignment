import http from "k6/http";
import { check } from "k6";

export let options = { vus: 50, iterations: 200 };

export default function () {
    const payload = JSON.stringify({
        user_id: "u3",
        is_active: false
    });

    const res = http.post("http://avito-test-assignment:8080/users/setIsActive", payload, {
        headers: { "Content-Type": "application/json" },
    });

    check(res, {
        "200 OK": r => r.status === 200,
    });
}
