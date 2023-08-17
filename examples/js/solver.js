const crypto = require('crypto');

function solveChallenge(data, criteria) {
    let i = 0;

    while (true) {
        const dataPlus = data + i;

        const hash = crypto.createHash('sha256').update(dataPlus).digest('hex');
        if (hash.slice(0, criteria.length) === criteria) {
            return {
                data: data,
                addedValue: i.toString(),
                hash: hash
            };
        }

        i++;
    }
}

// Example usage
const data = '1692065790855700|164|a5f7531fb66472a7a2524a9a5681e83ed389ed7db4c24895305d059719394b86';
const criteria = '000000';
console.log(solveChallenge(data, criteria));
