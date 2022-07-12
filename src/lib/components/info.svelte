<script lang="ts">
    import type { Stats } from "$lib/stats";


    export let stats: Stats;

    function prettyNumber(x: number) {
        return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
    }
</script>

<div class="container">
    <div>
        <h2>Downloads</h2>
        <p>Count: {prettyNumber(stats.downloads)}</p>
        <p>Online Players: {prettyNumber(stats.onlinePlayers)}</p>
    
        <div class="buttons">
            <button><img src="icons/download.svg" alt="download"> Release {stats.version} [{stats.mcVersion}]</button>
            <button><img src="icons/download.svg" alt="download"> Dev Build {stats.devBuildVersion} [{stats.devBuildMcVersion}]</button>
        </div>
    </div>
</div>

<style>
    .buttons {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
        margin-top: 1rem;
    }

    button {
        color: var(--text-primary);
        min-width: 20rem;
    }
</style>