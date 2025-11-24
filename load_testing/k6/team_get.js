import http from "k6/http";
import { check } from "k6";

export let options = { vus: 50, iterations: 200 };

export default function () {
    const res = http.get("http://avito-test-assignment:8080/team/get?team_name=backend");

    check(res, {
        "200 OK": r => r.status === 200,
    });
}
