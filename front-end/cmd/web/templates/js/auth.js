const output =document.getElementById("output");
const sent = document.getElementById("payload");
const received = document.getElementById("received");
const SubmitBtn = document.getElementById("submit");
const email = document.getElementById("email");
const password = document.getElementById("password");

SubmitBtn.addEventListener("click", function(event){
      event.preventDefault(); // prevent the form from submitting and reloading the page
      const inputemail = email.value;
      const inputpassword = password.value;

      const payload = {
        action : "auth",

        auth: {
          email: inputemail,
          password: inputpassword,

        }
      }
      const headers = new Headers();
      headers.append("Content-Type", "application/json");
      const body = {
         method: "POST",
         body : JSON.stringify(payload),
         headers: headers
      }
      fetch("http:\/\/localhost:8080/handle", body)
      .then((response) => response.json())
      .then((data) => {
          sent.innerHTML = JSON.stringify(payload, undefined, 4);
          received.innerHTML = JSON.stringify(data, undefined, 4);
          if (data.error){
             output.innerHTML += `<br><strong>Error:</strong>${data.message}`
          }else{
            output.innerHTML += `<br><strong>Response from broker service Email entered </strong>: ${data.email},
          ${data.message}`
          }
      })
      .catch(error => {
         output.innerHTML += "<br><br>Error: " + error;
      });
    })


