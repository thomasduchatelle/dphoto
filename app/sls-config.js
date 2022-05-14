module.exports = async ({options, resolveVariable}) => {
    const stage = await resolveVariable('sls:stage');
    // const region = await resolveVariable('opt:region, self:provider.region, "eu-west-1"');

    const rootDomain = `duchatelle.net`;
    const hostPrefix = stage === 'live' ? 'dphoto' : 'dphoto-dev'
    const domain = `${hostPrefix}.${rootDomain}`

    return {
        DPHOTO_JWT_KEY_B64: btoa(randomId(64)),
        DPHOTO_JWT_ISSUER: `https://${domain}`,
        DPHOTO_DOMAIN: domain,
        DPHOTO_ROOT_DOMAIN: rootDomain,
    }
}

function randomId(length) {
    let result = '';
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!£$%^&*(){}[]@~#|/.,<>?';
    const charactersLength = characters.length;
    for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() *
            charactersLength));
    }
    return result;
}