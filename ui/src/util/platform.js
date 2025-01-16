const platform = {
    isMobile: function() {
        return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    },
    isPortraitOrientation: function() {
        console.log(window.innerHeight, window.innerWidth);
        return (this.isMobile &&
            (window.screen.orientation.type === "portrait-primary" ||
            window.screen.orientation.type==="portrait-secondary")) || (window.innerWidth < window.innerHeight);
    },
}

export default platform;