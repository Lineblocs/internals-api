import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

export const errorRate = new Rate('non_success_requests');

export let options = {
  stages: [
    // Ramp-up from 1 to 10 VUs in 10s.
    { duration: "10s", target: 10 },

    // Stay at rest on 10 VUs for 5s.
    { duration: "5s", target: 10},

    //Linearly ramp down from 10 to 0 VUs over the last 15s.
    { duration: "15s", target: 0}
  ],
  thresholds: {
    // We want the 95th percentile of all HTTP request durations to be less than 500ms
    "http_req_duration": ["p(95)<1000"],
    // Thresholds based on the custom metric `non_success_requests`.
    "non_success_requests": [
      // Global failure rate should be less than 1%.
      "rate<0.01",
      // Abort the test early if it climb over 5%.
      { threshold: "rate<=0.05", abortOnFail: true},
    ],
  },
};

export default function () {

  const baseUrl = "http://localhost:8000";

  // Healthz
  /*
  const url = baseUrl + "/healthz";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */
  // CreateCall
  /*
  const url = baseUrl + "/call/createCall";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    from: "15874874526",
    to: "17808503688",
    status: "start",
    direction: "outbound",
    duration: "8",
    user_id: 2,
    workspace_id: 2,
    channel_id: "01gmvb228y61ay5b7afh9npxgx-ch",
    call_id: "1c11ae5423119c115143863e291d2f70@155.138.140.32:5160"
  });
  check(http.post(url,data, params), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // UpdateCall
  /*
  const url = baseUrl + "/call/updateCall";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    call_id: "call-b0d5d6f4-8ef4-4a8d-ac41-ffd7587c8eba",
    status: "ended"
  });
  check(http.post(url,data, params), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // FetchCall
  /*
  const url = baseUrl + "/call/fetchCall?id=256";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // SetSIPCallID
  /*
  const url = baseUrl + "/call/setSIPCallID";
  const data = {
    "callid": "aabc",
    "apiid": 255
  };
  check(http.post(url,data), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // SetProviderByIP
  /*
  const url = baseUrl + "/call/setProviderByIP";
  const data = {
    "ip": "toronto.voip.ms2",
    "apiid": 2572
  };
  check(http.post(url,data), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // CreateConference
  /*
  const url = baseUrl + "/conference/createConference";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    name: "test_conference4",
    workspace_id: 2
  });
  check(http.post(url,data,params), {
    'status is 200': (r) => r.status == 200,
  }) || errorRate.add(1);
  */

  // CreateDebit
  /*
  const url = baseUrl + "/debit/createDebit";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    user_id: 4,
    workspace_id: 2,
    module_id: 23,
    source: "aa",
    number: "40",
    type: "in",
    seconds: 240
  });
  check(http.post(url,data,params), {
    'status is 200': (r) => r.status == 200,
  }) || errorRate.add(1);
  */

  // CreateAPIUsageDebit
  /*
  const url = baseUrl + "/debit/createAPIUsageDebit";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    user_id: 4,
    workspace_id: 2,
    type: "STT",
    source: "test",
    params: {
        "length": 120,
        "recording_length": 50.5
    }
  });
  check(http.post(url,data,params), {
    'status is 200': (r) => r.status == 200,
  }) || errorRate.add(1);
  */

  // CreateLog
  /*
  const url = baseUrl + "/debugger/createLog";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    user_id: 2,
    workspace_id: 2,
    title: "Call Test",
    report: "Your call is connected successully",
    flow_id: 2,
    level: "info",
    from: "Me",
    to: "You"
  });
  check(http.post(url,data,params), {
    'status is 200': (r) => r.status == 200,
  }) || errorRate.add(1);
  */

  // CreateLogSimple
  /*
  const url = baseUrl + "/debugger/createLogSimple";
  const data = {
    "type": "verify-callerid-cailed",
    "level": "info",
    "domain": "workspace.com"
  };
  check(http.post(url,data), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // CreateFax
  /*
  const url = baseUrl + "/fax/createFax";
  const data = {
    "user_id": 2,
    "workspace_id": 2,
    "call_id": 255,
    "name": "test"
  };
  check(http.post(url,data), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // CreateRecording
  /*
  const url = baseUrl + "/recording/createRecording";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    workspace_id: 2,
    user_id: 2,
    call_id: 1,
    storage_id: "test",
    storage_server_ip: "test",
    trim: true,
    tags: ["a", "b", "c"]
  });
  check(http.post(url,data,params), {
    'status was 200': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // UpdateRecording
  /*
  const url = baseUrl + "/recording/updateRecording";
  const data = {
    "status": "completed",
    "recording_id": 5
  };
  check(http.post(url,data), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // UpdateRecordingTranscription
  /*
  const url = baseUrl + "/recording/updateRecordingTranscription";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    recording_id: 4,
    ready: false,
    text: "test2"
  });
  check(http.post(url,data,params), {
    'status was 200': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // GetRecording
  /*
  const url = baseUrl + "/recording/getRecording?id=4";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // CreateSIPReport
  /*
  const url = baseUrl + "/carrier/createSIPReport";
  const data = {
    "callid": 2,
    "status": 200
  };
  check(http.post(url,data), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // ProcessRouterFlow
  /*
  const url = baseUrl + "/carrier/processRouterFlow?callto=15874874526&callfrom=17808503688&userid=2";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // VerifyCaller
  /*
  const url = baseUrl + "/user/verifyCaller?workspace_id=2&number=15874874526";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // VerifyCallerByDomain
  /*
  const url = baseUrl + "/user/verifyCallerByDomain?domain=workspace&number=17808503688";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetUserByDomain
  /*
  const url = baseUrl + "/user/getUserByDomain?domain=workspace";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetUserByDID
  /*
  const url = baseUrl + "/user/getUserByDID?did=23";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetUserByTrunkSourceIp
  /*
  const url = baseUrl + "/user/getUserByTrunkSourceIp?source_ip=155.138.144.230";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetWorkspaceMacros
  /*
  const url = baseUrl + "/user/getWorkspaceMacros?workspace=2";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetDIDNumberData
  /*
  const url = baseUrl + "/user/getDIDNumberData?number=15874874526";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */
  
  // GetPSTNProviderIP
  /*
  const url = baseUrl + "/user/getPSTNProviderIP?from=15874874526&to=17808503688&domain=workspace";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetPSTNProviderIPForTrunk
  /*
  const url = baseUrl + "/user/getPSTNProviderIPForTrunk?from=15874874526&to=17808503688";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // IpWhitelistLookup
  /*
  const url = baseUrl + "/user/ipWhitelistLookup?ip=toronto.voip.ms&domain=workspace";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetDIDAcceptOption
  /*
  const url = baseUrl + "/user/getDIDAcceptOption?did=23";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetDIDAssignedIP
  /*
  const url = baseUrl + "/user/getDIDAssignedIP";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetUserAssignedIP
  /*
  const url = baseUrl + "/user/getUserAssignedIP?rtcOptimized=true&domain=workspace&routerip=155.138.140.32";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetTrunkAssignedIP
  /*
  const url = baseUrl + "/user/getTrunkAssignedIP";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // AddPSTNProviderTechPrefix
  /*
  const url = baseUrl + "/user/addPSTNProviderTechPrefix";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetCallerIdToUse
  /*
  const url = baseUrl + "/user/getCallerIdToUse?domain=workspace&extension=3000";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetExtensionFlowInfo
  /*
  const url = baseUrl + "/user/getExtensionFlowInfo?extension=3000&workspace=2";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetFlowInfo
  /*
  const url = baseUrl + "/user/getFlowInfo?flow_id=f-b8320fa4-68b2-4f5b-8a4a-dbcff9afb047&workspace=2";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetDIDDomain
  /*
  const url = baseUrl + "/user/getDIDDomain";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // GetCodeFlowInfo
  /*
  const url = baseUrl + "/user/getCodeFlowInfo?code=2&workspace=2";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // IncomingDIDValidation
  /*
  const url = baseUrl + "/user/incomingDIDValidation?did=23&number=15874874526&source=158.85.70.148";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // IncomingTrunkValidation
  /*
  const url = baseUrl + "/user/incomingTrunkValidation?fromdomain=155.138.144.230";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // LookupSIPTrunkByDID
  /*
  const url = baseUrl + "/user/lookupSIPTrunkByDID?did=23";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // IncomingMediaServerValidation
  /*
  const url = baseUrl + "/user/incomingMediaServerValidation?source=155.138.140.32";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // StoreRegistration
  /*
  const url = baseUrl + "/user/storeRegistration";
  const data = {
    "domain": "workspace",
    "user": 3000
  };
  check(http.post(url,data), {
    'status was success': (r) => r.status >= 200 && r.status <= 226,
  }) || errorRate.add(1);
  */

  // GetSettings
  /*
  const url = baseUrl + "/user/getSettings";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // ProcessSIPTrunkCall
  /*
  const url = baseUrl + "/user/processSIPTrunkCall?did=23";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  // SendAdminEmail
  /*
  const url = baseUrl + "/admin/sendAdminEmail";
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const data = JSON.stringify({
    message: "Hello. This is test message."
  });
  check(http.post(url,data,params), {
    'status is 200': (r) => r.status == 200,
  }) || errorRate.add(1);
  */

  // GetBestRTPProxy
  /*
  const url = baseUrl + "/getBestRTPProxy";
  let res = http.get(url);
  check(res, { "status was success": (r) => r.status >= 200 && r.status <= 226}) || errorRate.add(1);
  */

  sleep(Math.random()*1+1); //Random sleep betweeen 1s and 2s.
}