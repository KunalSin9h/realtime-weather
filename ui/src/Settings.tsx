import { useEffect, useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useQuery } from "@tanstack/react-query";

export default function Settings() {
  const [temperatureUnit, setTemperatureUnit] = useState('celsius');
  const [fetchInterval, setFetchInterval] = useState(3);


  const { data, isLoading, error } = useQuery({
    queryKey: ['user_preference'],
    queryFn: async () => {
      const resp = await fetch("http://localhost:5000/preference");
      if (!resp.ok || resp.status !== 200) {
        throw new Error("API return with non 200")
      }

      return resp.json();
    }
  })

  useEffect(() => {
    if (data) {
      setTemperatureUnit(data.time_unit);
      setFetchInterval(+data.interval.replaceAll("m0s", ""))
    }
  }, [isLoading])

  if (isLoading) {
    return <p className="text-gray-400">Loading...</p>
  }

  if (error) {
    return <p className="text-gray-400">Got Error: {error.message}</p>
  }

  return (
    <div className="p-4 max-w-md">
      <h2 className="text-2xl mb-4">Settings</h2>
      <div className="mb-4">
        <label htmlFor="temperatureUnit">
          Temperature Unit
        </label>
        <Select value={temperatureUnit} onValueChange={setTemperatureUnit}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Select unit" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="celsius">Celsius</SelectItem>
            <SelectItem value="kelvin">Kelvin</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="mb-4">
        <label htmlFor="fetchInterval">
          Fetch Interval <span>in minutes</span>
        </label>
        <Input
          type="number"
          id="fetchInterval"
          value={fetchInterval}
          min={3}
          max={60}
          inputMode="numeric"
          onChange={(e) => setFetchInterval(+e.target.value)}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
      </div>
      <Button onClick={() => {
        
      }}>
        Save
      </Button>
    </div>
  );
}
