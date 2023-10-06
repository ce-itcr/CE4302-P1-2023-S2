import React from "react";

export default function CardStats({
  statSubtitle,
  statTitle,
  code,
  lines,
  status,
  onClick,
  id
}) {
  return (
    <>
      <div className="relative flex flex-col min-w-0 break-words bg-white rounded mb-6 xl:mb-0 shadow-lg">
        <div className="flex-auto p-4">
          <div className="flex flex-wrap">
            <div className="relative w-full pr-4 max-w-full flex-grow flex-1">
              <span className="font-semibold text-xl text-blueGray-700">
                {statTitle} - reg {statSubtitle} 
              </span>
            </div>
            <div className="relative w-auto pl-4 flex-initial">
              <button
                className="bg-indigo-500 text-white active:bg-indigo-600 text-xs font-bold uppercase px-3 py-1 rounded outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
                type="button"
                style={{ backgroundColor: "#271744" }}
                onClick={onClick}
              >
                <i className="fas fa-forward-step text-white mr-4"></i>
                Step
              </button>

            </div>

          </div>
          <div className="flex flex-wrap" style={{ paddingTop: 10 }}>
            <div className="" >
              <label
                className="block uppercase text-blueGray-600 text-xs font-bold mb-2"
                htmlFor="full-name"
              >
                Código actual
              </label>
              <div style={{ height: 150, width: 190, backgroundColor: '#282c34', color: '#fff', padding: 20, borderRadius: 5, fontFamily: 'monospace', overflow: 'auto' }}>
                <span>
                  {code.map((line) => <p>{line.Instruction}</p>)}
                </span>
              </div>

            </div>
            <div className="" style={{ marginLeft: 60 }}>
              <label
                className="block uppercase text-blueGray-600 text-xs font-bold mb-2"
                htmlFor="full-name"
              >
                Lineas de caché
              </label>
              <div >
                <div >
                  <span>
                    <div className="flex flex-wrap">
                      <p style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', width: 30, height: 30, paddingLeft: 5 }}></p>
                      <p style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', width: 70, height: 30, paddingLeft: 5 }}>Bloque</p>
                      <p style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', width: 80, height: 30, paddingLeft: 5 }}>Dirección</p>
                      <p style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', width: 70, height: 30, paddingLeft: 5 }}>Dato</p>
                    </div>
                    {lines.map((line) =>
                      <div className="flex flex-wrap">
                        <p style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', width: 30, height: 30, paddingLeft: 5 }}>{line.State} </p>
                        <p style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', width: 70, height: 30, paddingLeft: 5 }}>{line.Block}</p>
                        <p style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', width: 80, height: 30, paddingLeft: 5 }}>{line.Address} </p>
                        <p style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', width: 70, height: 30, paddingLeft: 5 }}>{line.Data}</p>

                      </div>
                    )}
                  </span>
                </div>
              </div>

            </div>
          </div>
          <p className="text-sm " style={{ paddingTop: 10 }} >
            <span className={" mr-2"}>
              <i
                className={
                  "fas fa-arrow-right"
                }
              ></i>{" "}
              Estado:
            </span>
            <span className="whitespace-nowrap" style={{ color: '#40B714', fontWeight: 'bold' }}>{status}</span>
          </p>


        </div>

      </div>
    </>
  );
}
