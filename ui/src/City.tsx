import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import {useParams} from "react-router-dom";
import { Button } from "./components/ui/button";
import toast from "react-hot-toast";
import { LineChart } from "@/components/LineChart"

  
type LiveData = {
  city_id: number;
  time: string;
  temperature: number;
  humidity: number; 
  wind_speed: number;
}

export default function City() {
    const { city_name } = useParams();
    const { city_id } = useParams();
    const [noSummary, setNoSummary] = useState(false);
    const [liveData, setLiveData] = useState<LiveData[]>([]);
    const [lastUpdated, setLastUpdated] = useState<Date | null>(null);

    if (!localStorage.getItem("temp_unit")) {
        localStorage.setItem("temp_unit", "celsius");
    }

    const { data, isLoading, refetch } = useQuery({
        queryKey: ['cities_summary'],
        queryFn: async () => {
          const resp = await fetch(`/api/cities/summary/${city_id}`);
          if (!resp.ok) {
            return null;
          }
    
          // no contnet
          if (resp.status === 204) {
            setNoSummary(true);
          }
          return resp.json();
        }
      })

      useEffect(() => {

        const events = new EventSource(`/api/cities/live/${city_id}`)

        events.onmessage = (e) => {
          const data = JSON.parse(e.data);
          data.time = new Date(data.time).toLocaleTimeString();
          console.log(data);
          setLiveData((ld) => [...ld, data]);
          setLastUpdated(new Date());
        }
  
      }, [])

    
    if (isLoading) {
        return <div className="p-4">
            <p className="text-gray-400">Loading...</p>
        </div>
    }

    if (noSummary) {
        return <div className="p-4">
            <p className="text-xl font-bold mb-4">{city_name}</p>
            <p className="text-gray-600">No Weather Summary: Data is not yet processed</p>
            <p className="text-xs font-gray-400">Reason: You are seeing data before the schedule processing interval</p>
            <RefreshButton refetch={refetch} title="Calculate Now!" variant="default"/>
        </div>
    }

    return (
        <div className="p-4">
            <h2 className="text-3xl mb-4">{city_name}</h2>
            <div className="flex items-center gap mb-4">
                <div className="mr-8"> 
                    <p className="text-lg font-medium">Daily Weather Summary</p>
                    <p className="text-gray-400 text-xs">Updated Every 1 Hour.</p>
                </div>
                <RefreshButton variant="outline" refetch={refetch} title="Forcefully Recalculate"/>
            </div>
            <div className="grid">
              <dl className="bg-gray-100 rounded-md p-4 space-y-4">
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Date</p>
                  <p className="text-gray-500">Today</p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Avg Temperature</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.avg_temperature}</span>{' '}
                    °{localStorage.getItem('temp_unit')}
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Max Temperature</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.max_temperature}</span>{' '}
                    °{localStorage.getItem('temp_unit')}
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Min Temperature</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.min_temperature}</span>{' '}
                    °{localStorage.getItem('temp_unit')}
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Avg Humidity</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.avg_humidity}</span>%
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Max Humidity</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.max_humidity}</span>%
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Min Humidity</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.min_humidity}</span>%
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Avg Wind Speed</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.avg_wind_speed}</span>{' '}
                    km/h
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Max Wind Speed</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.max_wind_speed}</span>{' '}
                    km/h
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">Min Wind Speed</p>
                  <p className="text-gray-500">
                    <span className="text-blue-500">{data.min_wind_speed}</span>{' '}
                    km/h
                  </p>
                </dt>
                <dt className="flex items-center gap-4">
                  <p className="text-gray-700 text-sm uppercase">
                    Dominant Condition
                  </p>
                  <p className="text-blue-500">{data.dominant_condition}</p>
                </dt>
              </dl>
        </div>
        <div className="mt-8">
                <p>Weater Data</p>
                <p className="text-sm text-gray-400">Last Updated: {lastUpdated ? (
                  lastUpdated.getMinutes() > 2 ? (
                    lastUpdated.getMinutes() > 3 ? (
                      'Just Now'
                    ) : (
                      '2 minutes ago'
                    )
                  ) : (
                    '3 minutes ago'
                  )
                ) : 'N/A'}</p>
                <LineChart
    className="h-96 w-[60%]"
    data={liveData}
    index="time"
    categories={["temperature", "humidity", "wind_speed"]}
    valueFormatter={(number: number) => `${number}`}
    onValueChange={(v) => console.log(v)}
    xAxisLabel="Time"
    yAxisLabel="Data"
  />

              </div>
            </div>
    );
    return 
}

// @ts-expect-error
function RefreshButton({refetch, title, variant}) {
    const [isLoadingSummary, setIsLoadingSummary] = useState(false);

    return <Button variant={variant} className="mt-2" onClick={(e) => {
        e.preventDefault();
        setIsLoadingSummary(true);
    fetch(`/api/cities/summary/refresh`, {
      method: "POST",
    })
    .then((res) => {
      if (res.ok) {
        toast.success("Weather Summary is refreshed!");
      } else {
        toast.error("Failed to calculate weather summary");
      }
      setIsLoadingSummary(false);
      refetch();
    })
    .catch(() => {
      toast.error("Failed to calculate weather summary");
      setIsLoadingSummary(false);
      refetch();
    });
    }}>
        {isLoadingSummary ? "Doing!..." : title}
    </Button>
}