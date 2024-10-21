import { useState } from "react";

export default function Settings() {
  const [temperatureUnit, setTemperatureUnit] = useState('Celsius');
  const [fetchInterval, setFetchInterval] = useState(10);

  const handleTemperatureUnitChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setTemperatureUnit(event.target.value);
  };

  const handleFetchIntervalChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFetchInterval(parseInt(event.target.value, 10));
  };

  return (
    <div className="p-4 max-w-md mx-auto">
      <h2 className="text-2xl font-bold mb-4">Settings</h2>
      <div className="mb-4">
        <label htmlFor="temperatureUnit" className="block text-gray-700 font-bold mb-2">
          Temperature Unit
        </label>
        <select
          id="temperatureUnit"
          value={temperatureUnit}
          onChange={handleTemperatureUnitChange}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        >
          <option value="Celsius">Celsius</option>
          <option value="Kelvin">Kelvin</option>
        </select>
      </div>
      <div className="mb-4">
        <label htmlFor="fetchInterval" className="block text-gray-700 font-bold mb-2">
          Fetch Interval (minutes)
        </label>
        <input
          type="number"
          id="fetchInterval"
          value={fetchInterval}
          onChange={handleFetchIntervalChange}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
      </div>
      <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
        Save
      </button>
    </div>
  );
}
