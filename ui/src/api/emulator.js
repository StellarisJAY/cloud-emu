import api from "./request";

const emulator = {
    listEmulator: async function() {
        return await api.get("/emulator");
    },
    listGame: async function(emulatorType) {
        return await api.get("/game", {emulatorType: emulatorType, page: 1, pageSize: 1024});
    }
};

export default emulator;