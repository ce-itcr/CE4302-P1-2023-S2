import React, { useState } from "react";
import Header from "src/components/Header/Header";
import { useHistory } from "react-router-dom";
import Modal from "react-modal";
import { postRequest } from "src/common/communication";
import toast, { Toaster } from "react-hot-toast";

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

const Selector = () => {
  const protocols = ['MESI', 'MOESI']
  const [selectedOption, setSelectedOption] = useState(protocols[0]);
  const [lastCode, setLastCode] = useState(false);

  const [modalOpen, setModalOpen] = useState(false);

  const openModal = () => {
    setModalOpen(true);
  };
  const closeModal = () => {
    setModalOpen(false);
  };

  let history = useHistory();

  const onChangeLastCode = () => {
    setLastCode(!lastCode);
  }

  const startProcess = async () => {
    // console.log(selectedOption, !lastCode)
    localStorage.setItem('protocol', selectedOption);
    postRequest('setinitialize', { "type": selectedOption, "lastCode": !lastCode }).then((data) => {
      if (typeof data === 'undefined') {
        toast.error("Error de conexión");
        closeModal();
      } else {
        history.push('/app/dashboard');
      }
    })
  }

  return (
    <>
      <div><Toaster /></div>
      <Header title='/app/selector' />
      <section className="header relative pt-16 items-center flex h-screen max-h-860-px">

        <div className="container mx-auto px-4">
          <div className="items-center flex flex-wrap">


            <div className="w-full md:w-8/12 lg:w-6/12 xl:w-6/12 px-4">
              <div className="pt-32 sm:pt-0">
                <h2 className="font-semibold text-4xl text-blueGray-600">
                  Cache Sync - SW Cache Coherence modelling and evaluation.
                </h2>
                <p className="mt-4 text-lg leading-relaxed text-blueGray-500">
                  Cache Sync aims to analyze the performance of the MOESI and MESI cache coherence models in various situations, focusing on particular instructions (READ, WRITE, INC).
                </p>
                <p className="mt-4 text-lg leading-relaxed text-blueGray-500">
                  For detailed information on these specifications, we invite you to consult the {" "}
                  <a
                    href=""
                    className="text-blueGray-600"
                    target="_blank"
                    style={{ color: "#271744" }}
                  >
                    documentation
                  </a>
                  .
                </p>
              </div>
            </div>
            <div className="w-full ">

              <div className="absolute b-auto right-0  sm:w-6/12 -mt-48 sm:mt-0 w-10/12 max-h-860px">
                <div className="container mx-auto px-4">
                  <div className="flex flex-wrap justify-center lg:-mt-64 -mt-48">
                    <div className="w-full lg:w-6/12 px-4">
                      <div className="relative flex flex-col min-w-0 break-words w-full mb-6 shadow-lg rounded-lg bg-blueGray-200">
                        <div className="flex-auto lg:p-10">
                          {/* <h4 className="text-2xl font-semibold">
                            Selección de protocolo
                          </h4> */}
                          {/* <p className="leading-relaxed mt-1 mb-4 text-blueGray-500">
                            Para continuar, porfavor seleccione uno de los siguientes protocolos
                          </p> */}
                          <div className="relative w-full mb-3 ">
                            <label
                              className="block uppercase text-blueGray-600 text-xs font-bold mb-2"
                              htmlFor="full-name"
                            >
                              Protocol
                            </label>
                            <select
                              className="border-0 px-3 py-3 placeholder-blueGray-300 text-blueGray-600 bg-white rounded text-sm shadow focus:outline-none focus:ring w-full ease-linear transition-all duration-150"
                              value={selectedOption}
                              onChange={e => setSelectedOption(e.target.value)}
                            >
                              {protocols.map(currentProtocol => (
                                <option key={currentProtocol.value} value={currentProtocol}>{currentProtocol}</option>
                              ))}

                            </select>

                          </div>

                          <div className="relative w-full mb-3">
                            <label
                              className="block uppercase text-blueGray-600 text-xs font-bold mb-2"
                              htmlFor="email"
                            >
                              Random code
                            </label>
                            <input
                              id="customCheckLogin"
                              type="checkbox"
                              className="form-checkbox border-0 rounded text-blueGray-700 ml-1 w-5 h-5 ease-linear transition-all duration-150"
                              onChange={onChangeLastCode}
                            />
                            <span className="ml-2 text-sm font-semibold text-blueGray-600">
                              Use previously generated code
                            </span>
                          </div>


                          <div className="text-center mt-6">
                            <button
                              className="bg-blueGray-800 text-white active:bg-blueGray-600 text-sm font-bold uppercase px-6 py-3 rounded shadow hover:shadow-lg outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
                              type="button"
                              style={{ backgroundColor: "#271744" }}
                              onClick={openModal}
                            >
                              Select
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
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
            Start execution process
          </div>
          <div style={{ paddingBottom: 30 }}>
            Are you sure you want to start the execution process for the protocol <b style={{ color: "#271744" }}>{selectedOption}</b>  and the latest code {lastCode ? '' : 'not'} selected?
          </div>
          <div className="text-center mt-6">
            <button
              className="bg-blueGray-800 text-white active:bg-blueGray-600 text-sm font-bold uppercase px-6 py-3 rounded shadow hover:shadow-lg outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
              type="button"
              style={{ borderColor: "#271744", borderWidth: 1, color: "#271744", backgroundColor: '#fff', width: 150 }}
              onClick={closeModal}
            >
              Cancel
            </button>
            <button
              className="bg-blueGray-800 text-white active:bg-blueGray-600 text-sm font-bold uppercase px-6 py-3 rounded shadow hover:shadow-lg outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
              type="button"
              style={{ borderColor: "#271744", borderWidth: 1, backgroundColor: "#271744", width: 150 }}
              onClick={startProcess}
            >
              Start
            </button>
          </div>
        </div>
      </Modal>

    </>
  )

};

export default Selector;
