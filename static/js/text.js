var headerBox = document.getElementById("headertext");
var headerText = headerBox.innerHTML
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function typeText()
{
    headerBox.innerHTML = "‏‏‎ ‎";
    await sleep(300);
    headerBox.innerHTML = "";
    for (var i = 0; i < headerText.length; i++) {
        headerBox.innerHTML += headerText.charAt(i);
        await sleep(300);
    }
}
typeText()
