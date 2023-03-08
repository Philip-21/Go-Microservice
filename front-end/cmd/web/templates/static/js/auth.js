let output =document.getElementById("output");
let sent = document.getElementById("payload");
let received = document.getElementById("received");
let brokerBtn = document.getElementById("brokerBtn");
let authBrokerBtn = document.getElementById("authBrokerBtn");
let logBtn = document.getElementById("logBtn");
let logGBtn = document.getElementById("logGBtn");
let logRBtn = document.getElementById("logRBtn");


authBrokerBtn.addEventListener("click", function() {//event listener listens or clicks and executes the function
        const payload ={
            action:"auth", //what we expecting on the backend
            //data we are parsing
            auth:{
                email :"admin@example.com",
                password :"verysecret",
            }
        }
        const headers = new Headers();
        headers.append("Content-Type", "application/json");
        //body rep the html <body> in ref to the DOM ele dispalyed on a page
        const body ={
            method :"POST",
            body: JSON.stringify(payload),
            headers: headers,
        }
        //request to the server and load the info on the web page
        fetch("http:\/\/localhost:8080/handle", body)
        .then((response) => response.json()) //convert response to json
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