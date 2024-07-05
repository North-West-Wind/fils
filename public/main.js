
const button = document.querySelector("button#submit");
const resultDiv = document.querySelector("div#result");

button.onclick = async () => {
	if (button.classList.contains("disabled")) return;
	button.classList.add("disabled");
	resultDiv.classList.remove("ok", "err");
	resultDiv.style.display = "none";
	const link = document.querySelector("input#link").value;
	const pass = document.querySelector("input#pass").value;

	const res = await fetch("/api/new", {
		method: "POST",
		body: JSON.stringify({ link }),
		headers: {
			"Content-Type": "application/json",
			"Authorization": "Bearer " + pass
		}
	});
	resultDiv.style.display = "";
	if (res.ok) {
		resultDiv.classList.add("ok");
		resultDiv.innerHTML = `${window.location.href}s/${await res.text()}`;
	} else {
		resultDiv.classList.add("err");
		resultDiv.innerHTML = `Error: ${await res.text()}`;
	}

	button.classList.remove("disabled");
}

// copy to clipboard on click
let tmp, timeout;
resultDiv.onclick = () => {
	if (tmp) resultDiv.innerHTML = tmp;
  const buffer = document.createElement('textarea');
	buffer.value = resultDiv.innerHTML;
	resultDiv.appendChild(buffer);
	buffer.select();
  buffer.setSelectionRange(0, 99999);
  document.execCommand('copy');
  resultDiv.removeChild(buffer);

	tmp = resultDiv.innerHTML;
	resultDiv.innerHTML = "Copied!";
	if (timeout) clearTimeout(timeout);
	timeout = setTimeout(() => {
		if (tmp) {
			resultDiv.innerHTML = tmp;
			tmp = undefined;
		}
	}, 1000);
}