import {commonAPIGet} from "./common-http.js"

document.addEventListener("DOMContentLoaded", async function () {
    let urlSrc = "" 
    try {
        window.scrollTo(0, document.body.scrollHeight);
        const resp = await commonAPIGet("http://raccoon-mh.me:3000/api/pexelsrandomraccon")
        urlSrc = resp.data.responseData.photos[0].src.original
    } catch (error) {
        console.log(error)
        urlSrc = "https://images.pexels.com/photos/54602/raccoon-bear-zoo-saeugentier-54602.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1"
    }
    document.getElementById('cover-container').style.backgroundImage = 'url("'+urlSrc+'")'
    document.getElementById('cover-container').style.backgroundSize  = 'cover';
    document.getElementById('cover-container').style.backgroundPosition  = 'center';
});