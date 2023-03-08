
let output =document.getElementById("output");
let sent = document.getElementById("payload");
let received = document.getElementById("received");
let logBtn = document.getElementById("logBtn");


logBtn.addEventListener("click", function() {
      const payload = {
          action: "log",
          log: { //the log entry
              name: "event",
              data: "Some kind of data",
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
      fetch("http:\/\/localhost:8080/handle", body)
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