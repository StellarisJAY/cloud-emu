import api from "./request";

const emulator = {
    listEmulator: async function() {
        return await api.get("/emulator");
    },
    listGame: async function(emulatorId) {
        return await api.get("/game", {emulatorId: emulatorId, page: 1, pageSize: 1024});
    }
};

export default emulator;