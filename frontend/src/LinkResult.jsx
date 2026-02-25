import axios from "axios";
import { useEffect, useState } from "react"

const LinkResult = ({ inputValue }) => {
  const [shortenLink, setShortenLink] = useState("");
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);

  const fetchData = async () => {
    try {
      setLoading(true);
      // yaha kya likho ke mera mera api trigger 
      // http://localhost:3000/shorten
      const res = await axios(`http://localhost:3000/shorten?url=${inputValue}`);
      console.log("API response:", res);
      setShortenLink(res.data.short_url);
    } catch(err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    console.log("Input value changed:", inputValue);
    if(inputValue.length) {
      fetchData();
    }
  }, [inputValue]);

  useEffect(() => {
    const timer = setTimeout(() => {
      setCopied(false);
    }, 1000);

    return () => clearTimeout(timer);
  }, [copied]);

  if(loading) {
    return <p className="noData">Loading...</p>
  }
  if(error) {
    return <p className="noData">Something went wrong :(</p>
  }


  return (
    <>
      {shortenLink && (
        <div className="result">
          <p>{shortenLink}</p>
          <button
            className={copied ? "copied" : ""}
            onClick={async () => {
              try {
                await navigator.clipboard.writeText(shortenLink);
                setCopied(true);
              } catch (err) {
                console.error("Failed to copy:", err);
              }
            }}
          >
          {copied ? "Copied!" : "Copy to Clipboard"}
          </button>
        </div>
      )}
    </>
  )
}

export default LinkResult