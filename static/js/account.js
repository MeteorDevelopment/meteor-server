const token = localStorage.getItem("token")

function logout() {
    localStorage.removeItem("token")
    location.replace("/login")
}

if (token === null) {
    location.replace("/login")
}
else {
    document.addEventListener("DOMContentLoaded", () => {
        httpGetAuth("/api/account/info")
            .then(res => {
                console.log(res)

                populateDiscord(res)
                populateMcAccounts(res)
                populateCapes(res)
            })
            .catch(() => logout())
    })
}

// Discord

function populateDiscord(res) {
    $("#discord-preview").appendChild(createDiscordUsernameElement(res))

    let container = $("#discord-container")
    container.appendChild(createDiscordButtonElement(res.discord_id === ""))
}

function createDiscordUsernameElement(res) {
    let p = document.createElement("p")

    if (res.discord_id === "") {
        p.append("No Discord account linked.")
    }
    else {
        let li = document.createElement("li")
        li.classList.add("discord-account")
        p.appendChild(li)

        let img = document.createElement("img")
        img.src = res.discord_avatar
        img.style.width = "2em"
        img.alt = "discord avatar"
        li.appendChild(img)

        let b = document.createElement("b")
        b.append(res.discord_name)
        li.appendChild(b)
    }

    return p
}

function createDiscordButtonElement(link) {
    let btn = document.createElement("input")
    btn.type = "button"

    if (link) {
        btn.onclick = linkDiscord
        btn.value = "Link Discord"
    }
    else {
        btn.onclick = unlinkDiscord
        btn.value = "Unlink Discord"
    }

    return btn
}

function linkDiscord() {
    httpGetAuth("/api/account/generateDiscordLinkToken").then(res => {
        let container = $("#discord-container")
        container.replaceChildren()

        let p = document.createElement("p")
        p.append("To link your Discord account dm Meteor Bot on discord this message: ")
        container.appendChild(p)

        let b = document.createElement("b")
        b.append(".link " + res.token)
        p.appendChild(document.createElement("br"))
        p.appendChild(b)

        p.appendChild(document.createElement("br"))
        p.append("The command will only be valid for next 5 minutes.")
    })
}

function unlinkDiscord() {
    httpPostAuth("/api/account/unlinkDiscord").then(() => {
        let preview = $("#discord-preview")
        preview.replaceChildren()
        preview.appendChild(createDiscordUsernameElement({discord_id: ""}))

        let container = $("#discord-container")
        container.replaceChildren()
        container.appendChild(createDiscordButtonElement(true))
    })
}

// Mc accounts

function populateMcAccounts(res) {
    let mcAccountList = $("#mc-accounts-list")

    for (let i in res.mc_accounts) {
        createMcAccountElement(mcAccountList, res.mc_accounts[i])
    }

    let form = $("#add-mc-account-form")
    form.onsubmit = addMcAccount
    if (res.mc_accounts.length >= res.max_mc_accounts) form.remove()
}

function createMcAccountElement(parent, uuid) {
    let li = document.createElement("li")
    li.classList.add("mc-account-item")

    let img = document.createElement("img")
    img.src = "https://mc-heads.net/head/" + uuid + "/32"
    img.alt = "head"
    li.appendChild(img)

    let btn = document.createElement("button")
    btn.classList.add("remove-mc-account")
    btn.onclick = () => httpDeleteAuth("/api/account/mcAccount?uuid=" + uuid).then(() => location.reload())
    btn.append("Remove")

    httpGet("https://mc-heads.net/minecraft/profile/" + uuid)
        .then(res => li.append(res.name))
        .catch(() => li.append(uuid))
        .finally(() => {
            li.appendChild(btn)
            parent.appendChild(li)
        })
}

function addMcAccount(e) {
    e.preventDefault()

    let name = $("#add-mc-account-username").value
    if (name !== "") {
        httpPostAuth("/api/account/mcAccount?username=" + name).then(() => location.reload())
    }
}

// Capes

function populateCapes(res) {
    let capeList = $("#cape-list")

    for (let i in res.capes) {
        capeList.appendChild(createCapeElement(res.capes[i]))
    }

    let form = $("#upload-cape-form")
    let error = $("#upload-cape-error")

    form.onsubmit = e => {
        e.preventDefault()

        let file = $("#upload-cape-file").files[0]
        let form = new FormData()
        form.append("file", file)

        error.replaceChildren()

        httpPostAuthBody("/api/account/uploadCape", form)
            .then(() => location.reload())
            .catch(res => {
                let p = document.createElement("p")
                p.classList.add("error")
                p.append(res.error)
                error.appendChild(p)
            })
    }

    if (!res.can_have_custom_cape) form.remove()
}

function createCapeElement(cape) {
    let li = document.createElement("li")
    li.classList.add("cape-item")

    if (cape.url !== "") {
        let img = document.createElement("img")
        img.src = cape.url
        img.alt = "cape preview"
        img.style.width = "10em"
        li.appendChild(img)
    }

    if (cape.current) {
        let strong = document.createElement("strong")
        strong.append(cape.title)
        li.appendChild(strong)
    }
    else {
        let p = document.createElement("p")
        p.id = "cape-title"
        p.append(cape.title)
        li.appendChild(p)

        let btn = document.createElement("button")
        btn.classList.add("select-cape")
        btn.append("Select")
        btn.onclick = () => httpPostAuth("/api/account/selectCape?cape=" + cape.id).then(() => location.reload())
        li.appendChild(btn)
    }

    return li
}

// Controls

function logoutBtn() {
    httpPostAuth("/api/account/logout").finally(() => logout())
}

function changeUsername() {
    location.replace("/changeUsername")
}

function changeEmail() {
    location.replace("/changeEmail")
}

function changePassword() {
    location.replace("/changePassword")
}
