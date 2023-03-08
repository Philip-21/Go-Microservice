let output =document.getElementById("output");
let sent = document.getElementById("payload");
let received = document.getElementById("received");
let brokerBtn = document.getElementById("brokerBtn");
let authBrokerBtn = document.getElementById("authBrokerBtn");
let logBtn = document.getElementById("logBtn");
let logGBtn = document.getElementById("logGBtn");
let logRBtn = document.getElementById("logRBtn");


logGBtn.addEventListener("click", function() {
      //Rest-Api Actions
      const payload = {
          action: "log",
          log: { 
              name: "Grpc Event",
              data: "Grpc Api Calls ",
          }
      }
      //Post headers
      const headers = new Headers();
      headers.append("Content-Type", "application/json");
      const body = {
          method: "POST",
          body: JSON.stringify(payload),
          headers: headers,
      }
      //calls the Grpc Handler to execute the request
      fetch("http:\/\/localhost:8080/log-grpc", body)
      .then((response) => response.json())
      .then((data) => {
          sent.innerHTML = JSON.stringify(payload, undefined, 4);
          received.innerHTML = JSON.stringify(data, undefined, 4);
          if (data.error) {
              output.innerHTML += `<br><strong>Error:</strong> ${data.message}`;
          } else {
              output.innerHTML += `<br><strong>Response from broker service</strong>: ${data.message}`;
          }
      })
      .catch((error) => {
          output.innerHTML += "<br><br>Eror: " + error;
      })        
  })
