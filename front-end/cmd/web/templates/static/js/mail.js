let output =document.getElementById("output");
let sent = document.getElementById("payload");
let received = document.getElementById("received");
let mailBtn = document.getElementById("mailBtn");




        mailBtn.addEventListener("click", function() {
            const payload = {
                action: "mail",
                mail: {
                    from: "me@example.com",
                    to: "you@there.com",
                    subject: "Test email",
                    message: "Hello world!",
                }
            }
            const headers = new Headers();
            headers.append("Content-Type", "application/json");
            const body = {
                method: 'POST',
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