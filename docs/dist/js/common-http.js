export async function commonAPIPost(url, data) {
    const response = await axios.post(url, data)
        .then(function (response) {
            return response
        })
        .catch(function (error) {
            return error
        });

    return response
}

export async function commonAPIGet(url) {
    const response = await axios.get(url)
        .then(function (response) {
            return response
        })
        .catch(function (error) {
            return error
        });

    return response
}

export async function commonAPIPostDEBUG(url, data) {
    console.log("#### commonAPIPost")
    console.log("Request URL : ", url)
    console.log("Request Data : ")
    console.log(JSON.stringify(data))
    console.log("-----------------------")
    const response = await axios.post(url, data)
        .then(function (response) {
            console.log("#### commonAPIPost Response")
            console.log("Response status : ", (response.status))
            console.log("Response Data :")
            console.log(JSON.stringify(response.data))
            console.log("----------------------------")
            return response
        })
        .catch(function (error) {
            console.log("#### commonAPIPost Response ERR")
            console.log("error : ", (error))
            console.log("--------------------------------")
            return error
        });

    return response
}

export async function commonAPIGetDEBUG(url) {
    console.log("#### commonAPIGet")
    console.log("Request URL : ", url)

    const response = await axios.get(url)
        .then(function (response) {
            console.log("#### commonAPIPost Response")
            console.log("Response status : ", (response.status))
            console.log("Response Data :")
            console.log(response.data)
            console.log("----------------------------")
            return response
        })
        .catch(function (error) {
            console.log("#### commonAPIPost Response ERR")
            console.log("error : ", (error))
            console.log("--------------------------------")
            return error
        });

    return response
}