import http from "k6/http";
import { sleep, check } from "k6";

export const options = {
  discardResponseBodies: false,

  // SIMPLE, BRUTAL, AND EFFECTIVE: Constant Virtual Users
  scenarios: {
    constant_vus_siege: {
      executor: "constant-vus", // The simplest executor
      vus: 10_000, // MAX OUT YOUR LOCAL HARDWARE (2000 VUs)
      duration: "10m", // Siege for 10 minutes straight
    },
  },

  thresholds: {
    http_req_failed: ["rate<0.95"],
    http_req_duration: ["p(99)<10000"],
  },
};

export default function () {
  let res = http.get("http://localhost:8090/teams/?page=1&limit=50");
  let checkResult = check(res, {
    "status is 200": (res) => res.status === 200,
  });
  if (!checkResult) {
    fail(`Request failed with status ${res.status}`);
  }
  sleep(1);
}
