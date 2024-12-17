const tokenStorage = {
    getToken() {
        return sessionStorage.getItem("cloudemu_token")
    },
    setToken(token) {
        sessionStorage.setItem("cloudemu_token", token)
    },
    delToken() {
        sessionStorage.clear()
    }
}

export default tokenStorage