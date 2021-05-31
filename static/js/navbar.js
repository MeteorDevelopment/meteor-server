const navSlide = () => {
    const burger = $(".burger")
    const nav = $(".nav-links")
    const navLinks = document.querySelectorAll(".nav-links li")

    // Toggle borgor
    burger.addEventListener("click", () => {
        nav.classList.toggle("nav-active")

        // Animate
        navLinks.forEach((link, index) => {
            if (link.style.animation) {
                link.style.animation = ``
            }
            else {
                link.style.animation = `navLinkFade 0.5s ease forwards ${
                    index / 7 + 0.1
                }s`
            }
        })

        // Animate da borgor
        burger.classList.toggle("toggle");
    })
}

navSlide()