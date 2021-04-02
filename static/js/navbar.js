const burger = document.querySelector(".burger");
const list = document.querySelector("nav ul");
const items = document.querySelectorAll("nav ul li");

burger.addEventListener("click", () => {

    list.classList.toggle("active");

    items.forEach((link, index) => {
        if (link.style.animation) {
            link.style.animation = ``;
        } else {
            link.style.animation = `navLinkFade 0.5s ease forwards ${
                index / 7 + 0.1
            }s`;
        }
    });

    burger.classList.toggle("active");

});