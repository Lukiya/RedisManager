import { GetCluster } from "@/services/cluster";

const _defaultCluster = {
    ID: "selected",
    Nodes: [],
};

export default {
    state: {
        cluster: { _defaultCluster },
        openKeys: [],
    },
    effects: {
        *getCluster({ clusterID }: any, { call, put }: any): any {
            const resp = yield call(GetCluster, clusterID);
            if (!resp.Nodes) {
                yield put({
                    type: 'setState',
                    payload: {
                        cluster: _defaultCluster,
                        openKeys: [],
                    },
                });
                return;
            }
            let selectedNodeID = null;
            for (let i = 0; i < resp.Nodes.length; i++) {
                selectedNodeID = resp.Nodes[0].ID;
                break;
            }

            yield put({
                type: 'setState',
                payload: {
                    cluster: resp,
                    openKeys: [selectedNodeID],
                },
            });
        },
        *refresh({ clusterID }: any, { call, put, select }: any): any {
            const state = yield select((x: any) => x["menuVM"]);
            if (clusterID == "selected" || state.cluster.ID == clusterID) {    // only refresh selected cluster
                yield put({
                    type: 'getCluster',
                    clusterID,
                });
            }
        },
    },
    reducers: {
        setState(state: any, { payload }: any) { return { ...state, ...payload }; },
        setOpenKeys(state: any, { openKeys }: any) {
            if (openKeys.length > 0) {
                openKeys = [openKeys[openKeys.length - 1]];
            }
            return {
                ...state,
                openKeys,
            };
        },
    },
};