<script lang="ts">
    import type { User } from "$lib/user";

    export let user: User;
</script>

<nav>
<!--    <img src="icon.png" alt="icon">-->
    <a href="/">Home</a>
    <a href="/addons">Addons</a>
    <a href="/faq">FAQ</a>
    <a href="/github">GitHub</a>
    {#if user}
        <p>{user.username}</p>
        <img src={user.discord_avatar} alt="icon" />
    {:else}
        <button>Login</button>
    {/if}
</nav>

<style>
    nav {
        background-color: var(--accent);
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 2rem;
    }

    nav p {
        position: absolute;
        right: 4rem;
        direction: rtl;
        color: var(--text-primary);
    }

    nav img {
        position: absolute;
        right: 1rem;
        max-height: 40px;
        max-width: 40px;
        border-radius: 50%;
    }

    img {
        width: 4rem;
    }

    a, p {
        text-decoration: none;
        color: var(--text-secondary);
        padding-top: 1rem;
        padding-bottom: 1rem;;
    }
</style>