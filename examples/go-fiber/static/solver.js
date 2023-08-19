async function solveChallenge(data, criteria) {
    let i = 0;

    while (true) {
        const dataPlus = data + i;
        const encoder = new TextEncoder();
        const hashBuffer = await crypto.subtle.digest('SHA-256', encoder.encode(dataPlus));
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');

        if (hashHex.slice(0, criteria.length) === criteria) {
            return {
                data: data,
                addedValue: i.toString(),
                hash: hashHex
            };
        }

        i++;
    }
}
