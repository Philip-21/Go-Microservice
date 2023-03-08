let output =document.getElementById("output");
let sent = document.getElementById("payload");
let received = document.getElementById("received");
let brokerBtn = document.getElementById("brokerBtn");
let authBrokerBtn = document.getElementById("authBrokerBtn");
let logBtn = document.getElementById("logBtn");
let logGBtn = document.getElementById("logGBtn");
let logRBtn = document.getElementById("logRBtn");



logRBtn.addEventListener("click",function(){
      const payload ={
          action: "rpc",
          log :{
              name : "rpc entry",
              data: "rpc api calls",
          }
      }
      const headers = new Headers();
      headers.append("Content-Type","application/json");

      const body = {
          method: "POST",
          body: JSON.stringify(payload),
          headers: headers,
      }
      fetch("http:\/\/localhost:8080/handle", body)
      .then((response) => response.json())
      .then((data) => {
          sent.innerHTML =JSON.stringify(payload, undefined, 4);
          received.innerHTML =JSON.stringify(data, undefined, 4); //formatin JSON
          if (data.error) {
             output.innerHTML += `<br><strong>Error:</strong> ${data.message}`;
          }else{
              output.innerHTML += `<br><strong>Response from broker service</strong>: ${data.message}`;
          }
      })
      .catch((error) => {
          output.innerHTML += "<br><br>Eror: " + error;
      })
  })