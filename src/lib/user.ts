export type User = {
    id: string;
    username: string;
    email: string;
    admin: boolean;
    donator: boolean;
    discord_id: bigint;
    max_mc_accounts: number;
    mc_accounts: string[];
    cape: string;
    can_have_custom_cape: boolean;
    discord_name: string;
    discord_avatar: string;
    capes: string[];
}

export async function fetchUser(fetch: any, token: string): Promise<User> {
    const response = await fetch('https://meteorclient.com/api/account/info', {
        method: 'GET',
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });

    if (!response.ok) console.warn(`Failed to fetch account info. Status: ${response.status}`);

    const user = (await response.json());

    return {
        id: user.id,
        username: user.username,
        email: user.email,
        admin: user.admin,
        donator: user.donator,
        discord_id: user.discord_id,
        max_mc_accounts: user.max_mc_accounts,
        mc_accounts: user.mc_accounts,
        cape: user.cape,
        can_have_custom_cape: user.can_have_custom_cape,
        discord_name: user.discord_name,
        discord_avatar: user.discord_avatar,
        capes: user.capes
    };
}