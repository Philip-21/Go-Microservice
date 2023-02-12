const output =document.getElementById("output");
const sent = document.getElementById("payload");
const received = document.getElementById("received");
const SubmitBtn = document.getElementById("submit");
const firstname = document.getElementById("firstname");
const lastname = document.getElementById("lastname");
const email = document.getElementById("email");
const password = document.getElementById("password");



SubmitBtn.addEventListener("click", function(event){
      event.preventDefault;
      inputfirstname= firstname.value;
      inputlastname = lastname.value;
      inputemail = email.value;
      inputpassword = password.value

      const payload = {
            action: "signup",
            signup: {
                  firstname: inputfirstname,
                  lastname: inputlastname,
                  email : inputemail,
                  password: inputpassword
            }
      }
      const headers= new Headers();
      headers.append("Content-Type", "application/json");
      const body ={
            method : "POST",
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
            output.innerHTML += `<br><strong>Response from authentication to broker service signup Successfuly </strong>: ${data.firstname}`
          }
      })
      .catch(error => {
         output.innerHTML += "<br><br>Error: " + error;
      });
})