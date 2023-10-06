import React, { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { getMetricsData } from "src/common/communication";
import CardResume from "src/components/Cards/CardResume";
import CardTable from "src/components/Cards/CardTable";
import Header from "src/components/Header/Header";


const Summary = () => {

    const [transactions, setTransactions] = useState([]);
    const [metrics, setMetrics] = useState({});


    let history = useHistory();

    const checkOnClose = () => {
        history.push('/app');
    }

    const getData = () => {
        getMetricsData('').then((data) => {
            console.log(data)
            setTransactions(data.Transactions)
            console.log(data.Transactions)
            let currentMetrics = { "PowerConsumption": data.PowerConsumption.toFixed(2), "CacheMisses": data.CacheMisses, "MemoryAccesses": data.MemoryAccesses, "MissRate": data.MissRate, "HitRate": data.HitRate, "ReadRequests": data.ReadRequests, "ReadExclusiveRequest": data.ReadExclusiveRequest, "DataResponses": data.DataResponses, "Invalidates": data.Invalidates, "MemoryReads": data.MemoryReads, "MemoryWrites": data.MemoryWrites };
            console.log("🚀 ~ file: Summary.js:26 ~ getMetricsData ~ currentMetrics:", currentMetrics)
            setMetrics(currentMetrics);
        })
    }

    useEffect(() => {
        getData();

    }, [])



    return (
        <>
            <Header title={'/app/summary'} onClose={checkOnClose} onDownload={checkOnClose} />

            <section className="header relative pt-16 items-center flex h-screen max-h-860-px">

                <div className="container mx-auto px-4">
                    <div className="items-center flex flex-wrap">


                        <div className="w-full md:w-8/12 lg:w-6/12 xl:w-6/12 px-4">
                            <div className="pt-32 sm:pt-0">
                                <CardTable transactions={transactions} />

                            </div>
                        </div>

                        <div className="w-full md:w-8/12 lg:w-6/12 xl:w-6/12 px-4">
                            <div className="pt-32 sm:pt-0">
                                <CardResume metrics={metrics} />
                            </div>
                        </div>


                    </div>
                </div>

            </section>
        </>
    )

};

export default Summary;
