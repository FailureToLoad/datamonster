import { useParams } from "react-router-dom";
import { columns } from "./columns";
import { DataTable } from "./data-table";
import { useAuth0 } from "@auth0/auth0-react";
import { Survivor } from "./types";
import { Get } from "@/api/api";
import { useQuery } from "@tanstack/react-query";
import Spinner from "@/components/spinner";

export default function Population() {
  const { settlementId } = useParams();
  const { getAccessTokenSilently, isLoading } = useAuth0();
  const getPopulation = async () => {
    try {
      const token = await getAccessTokenSilently();
      const response = await Get<Array<Survivor>>(
        `settlement/${settlementId}/survivor`,
        token,
      );
      if (!response.data) return null;
      return response.data;
    } catch (e) {
      console.log(e);
      return null;
    }
  };
  const { isPending, isError, data, error } = useQuery({
    queryKey: ["population"],
    queryFn: getPopulation,
  });

  if (isPending || isLoading) {
    return <Spinner />;
  }

  if (isError) {
    throw new Error(error.message);
  }

  let population = data as Array<Survivor>;
  return (
    <div id="population">
      <DataTable columns={columns} data={population} />
    </div>
  );
}
