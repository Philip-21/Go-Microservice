const output =document.getElementById("output");
const sent = document.getElementById("payload");
const received = document.getElementById("received");
const SubmitBtn = document.getElementById("submit");
const email = document.getElementById("email");
const password = document.getElementById("password");

SubmitBtn.addEventListener("click", function(event){
  event.preventDefault(); // prevent the form from submitting and reloading the page
})
