const formData = new FormData();
formData.append("pwd", localStorage.getItem("pwd"));

const obj = {
  method: "POST",
  body: formData,
}

Util.fetch({
  url: "/api/check-pwd",
  obj: obj,
  onError: (err) => {
    // window.location.href = "login.html";
  }
});
