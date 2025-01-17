import api from "./request";

const emulator = {
    listEmulator: async function() {
        return await api.get("/emulator");
    },
    listGame: async function(emulatorType) {
        return await api.get("/game", {emulatorType: emulatorType, page: 1, pageSize: 1024});
    },
    listEmulatorTypes: async function() {
        return await api.get("/emulator/type");
    },
    updateEmulator: async function(emulator)  {
        return await api.put("/emulator", emulator);
    },
};

export default emulator;