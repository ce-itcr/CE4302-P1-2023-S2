import React, { useEffect, useState } from "react";
import CardLines from "src/components/Cards/CardLines";
import CardLogger from "src/components/Cards/CardLogger";
import CardStats from "src/components/Cards/CardStats";
import Header from "src/components/Header/Header";
import { useHistory } from "react-router-dom";
import Modal from "react-modal";
import data from '../../common/info.json';
import { sleep } from "src/utils/Sleep";
import { getDashboardData, postRequest } from "src/common/communication";
import { useInterval } from "src/utils/useInterval";


export const customStyles = {
    content: {
        backgroundColor: "#fff",
        color: "#000",
        top: "50%",
        left: "50%",
        right: "auto",
        bottom: "auto",
        marginRight: "-50%",
        transform: "translate(-50%, -50%)",
    },
};


const Dashboard = () => {
    const [currentProtocol, setCurrentProtocol] = useState('');
    const [ICLogs, setICLogs] = useState([]);
    const [memory, setMemory] = useState([]);

    const [register0, setRegister0] = useState([]);
    const [register1, setRegister1] = useState([]);
    const [register2, setRegister2] = useState([]);

    const [code0, setCode0] = useState([]);
    const [code1, setCode1] = useState([]);
    const [code2, setCode2] = useState([]);

    const [lines0, setLines0] = useState([]);
    const [lines1, setLines1] = useState([]);
    const [lines2, setLines2] = useState([]);

    const [status0, setStatus0] = useState('');
    const [status1, setStatus1] = useState('');
    const [status2, setStatus2] = useState('');


    const [dashboardData, setDashboardData] = useState({})


    let history = useHistory();

    const [modalOpen, setModalOpen] = useState(false);

    const openModal = () => {
        setModalOpen(true);
    };
    const closeModal = () => {
        setModalOpen(false);
    };

    const getData = () => {

        getDashboardData('').then((data) => {
            setDashboardData(data)

            setICLogs(data.IC.Logs)
            setMemory(data.MM.Blocks)

            const currentCode0 = data.PEs.find(pe => pe.ID === 0)?.Instructions;
            const currentRegister0 = data.PEs.find(pe => pe.ID === 0)?.Register;
            const currentStatus0 = data.PEs.find(pe => pe.ID === 0)?.Status;
            const currentLines0 = data.CCs.find(cc => cc.ID === 0)?.Cache;
            setCode0(currentCode0)
            setRegister0(currentRegister0)
            setStatus0(currentStatus0)
            setLines0(currentLines0)

            const currentCode1 = data.PEs.find(pe => pe.ID === 1)?.Instructions;
            const currentRegister1 = data.PEs.find(pe => pe.ID === 1)?.Register;
            const currentStatus1 = data.PEs.find(pe => pe.ID === 1)?.Status;
            const currentLines1 = data.CCs.find(cc => cc.ID === 1)?.Cache;
            setCode1(currentCode1)
            setRegister1(currentRegister1)
            setStatus1(currentStatus1)
            setLines1(currentLines1)

            const currentCode2 = data.PEs.find(pe => pe.ID === 2)?.Instructions;
            const currentRegister2 = data.PEs.find(pe => pe.ID === 2)?.Register;
            const currentStatus2 = data.PEs.find(pe => pe.ID === 2)?.Status;
            const currentLines2 = data.CCs.find(cc => cc.ID === 2)?.Cache;
            setCode2(currentCode2)
            setRegister2(currentRegister2)
            setStatus2(currentStatus2)
            setLines2(currentLines2)

        })

    }

    useEffect(() => {

        setCurrentProtocol(localStorage.getItem('protocol'));

        getData();

    }, [])

    const checkOnResume = async () => {
        history.push('/app/summary');

        postRequest('setlj', {
            "action": "close"
        }).then(() => {
            history.push('/app/summary');
        })

    }

    useInterval(() => {
        getData();
    }, 500)

    const onStep = async (pe) => {
        switch (pe) {
            case 1:
                postRequest('setaction', {
                    "action": "step",
                    "number": "0"
                })
                break;
            case 2:
                postRequest('setaction', {
                    "action": "step",
                    "number": "1"
                })
                break;
            case 3:
                postRequest('setaction', {
                    "action": "step",
                    "number": "2"
                })
                break;
            default:
                break;
        }
    }

    const onStart = async () => {
        postRequest('setaction', {
            "action": "all",
            "number": "1"
        })
    }

    // useEffect(() => {
    //     const intervalId = setInterval(() => {
    //       console.log('Acción realizada cada 3 segundos');
    //     }, 3000); 

    //     return () => {
    //       clearInterval(intervalId);
    //     };
    //   }, []); 



    return (
        <>
            <Header title={'/app/dashboard'} onResume={openModal} />

            <section className="header relative pt-16 items-center flex h-screen max-h-860-px">

                <div className="container mx-auto px-4">
                    <div className="items-center flex flex-wrap">


                        <div className="w-full md:w-8/12 lg:w-6/12 xl:w-6/12 px-4">
                            <div className="pt-32 sm:pt-0">
                                <div style={{ paddingTop: 60, paddingBottom: 5 }}>
                                    <CardStats
                                        statSubtitle={register0}
                                        statTitle="PE#1"
                                        id="0"
                                        code={code0}
                                        lines={lines0}
                                        status={status0}
                                        onClick={() => { onStep(1) }}
                                    />
                                </div>
                                <div style={{ paddingBottom: 5 }}>
                                    <CardStats
                                        statSubtitle={register1}
                                        statTitle="PE#2"
                                        id="1"
                                        code={code1}
                                        lines={lines1}
                                        status={status1}
                                        onClick={() => { onStep(2) }}
                                    />
                                </div>
                                <div>
                                    <CardStats
                                        statSubtitle={register2}
                                        statTitle="PE#3"
                                        id="2"
                                        code={code2}
                                        lines={lines2}
                                        status={status2}
                                        onClick={() => { onStep(3) }}
                                    />
                                </div>

                            </div>
                        </div>

                        <div className="w-full md:w-8/12 lg:w-6/12 xl:w-6/12 px-4">
                            <div className="pt-32 sm:pt-0">
                                <div style={{ paddingTop: 60, paddingBottom: 20 }}>
                                    <CardLines onClick={ onStart } data={memory} />
                                </div>
                                <div>
                                    <CardLogger logs={ICLogs} protocol={currentProtocol} />
                                </div>

                            </div>
                        </div>


                    </div>
                </div>

            </section>

            <Modal
                isOpen={modalOpen}
                onRequestClose={closeModal}
                style={customStyles}
            >
                <div style={{ maxWidth: 400 }}>
                    <div style={{ fontWeight: 'bold', fontSize: 18, paddingBottom: 20, color: "#271744" }}>
                        Detener proceso de ejecución
                    </div>
                    <div style={{ paddingBottom: 30 }}>
                        ¿Está seguro que desea detener el proceso de ejecución?
                    </div>
                    <div className="text-center mt-6">
                        <button
                            className="bg-blueGray-800 text-white active:bg-blueGray-600 text-sm font-bold uppercase px-6 py-3 rounded shadow hover:shadow-lg outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
                            type="button"
                            style={{ borderColor: "#271744", borderWidth: 1, color: "#271744", backgroundColor: '#fff', width: 150 }}
                            onClick={closeModal}
                        >
                            Cancelar
                        </button>
                        <button
                            className="bg-blueGray-800 text-white active:bg-blueGray-600 text-sm font-bold uppercase px-6 py-3 rounded shadow hover:shadow-lg outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
                            type="button"
                            style={{ borderColor: "#271744", borderWidth: 1, backgroundColor: "#271744", width: 150 }}
                            onClick={checkOnResume}
                        >
                            Detener
                        </button>
                    </div>
                </div>
            </Modal>

        </>
    )

};

export default Dashboard;
