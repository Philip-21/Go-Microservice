let output =document.getElementById("output");
let sent = document.getElementById("payload");
let received = document.getElementById("received");
let brokerBtn = document.getElementById("brokerBtn");



//a listener on the broker button that listens for clicks
brokerBtn.addEventListener("click", function(){
      const body = {
          //post an empty body
          method: 'POST',
      }
   
      fetch("http:\/\/localhost:8080", body) //a standard way of parsin a url in go templates
      .then((response) => response.json()) //convert response to json
      .then((data) => {
          sent.innerHTML ="empty post request";
          received.innerHTML =JSON.stringify(data, undefined, 4); //formatin JSON
          if (data.error) {//displays false error
              console.log(data.message); //displays the messae hit the broker
          }else{
              output.innerHTML += `<br><strong>Response from broker service</strong>: ${data.message}`; //message hit the broker 
          }
      })
      .catch((error) => {
          output.innerHTML += "<br><br>Eror: " + error;
      })
  })
