<script lang="ts">
    import Navbar from "$lib/components/navbar.svelte";
    import Main from "$lib/components/main.svelte";
    import Info from "$lib/components/info.svelte";
    import Mpvp from "$lib/components/mpvp.svelte";
    import Team from "$lib/components/team.svelte";
    import Footer from "$lib/components/footer.svelte";
</script>

<script lang="ts" context="module">
    import type { Load } from "@sveltejs/kit";
    import { fetchStats } from "$lib/stats";
    import { fetchUser } from "$lib/user";
    // import { writable } from 'svelte/store';
    // import type { Writable } from "svelte/types/runtime/store";
    //
    // export const tokenStore: Writable<string> = writable("");
    // let token = "";
    //
    // tokenStore.subscribe(value => token = value);

    export const load: Load = async ({ fetch }) => {
        return {
            props: {
                stats: await fetchStats(fetch),
                user: await fetchUser(fetch, null /*TODO*/)
            }
        }
    };
</script>

<Navbar user={$$props.user}/>
<Main />
<Info stats={$$props.stats} />
<Mpvp />
<Team />
<Footer />