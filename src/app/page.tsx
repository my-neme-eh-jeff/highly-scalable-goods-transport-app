import {
  BackgroundShape1,
  BackgroundShape2,
  BackgroundShape3,
} from "@/assets/svg";
import { ArrowRight, Car, DollarSign, MapPinned } from "lucide-react";
import Link from "next/link";

export default function Home() {
  const negText = {
    color: "#F26602",
    mixBlendMode: "difference",
  };

  const negDiv = {
    zIndex: 10,
    backgroundColor: "#F26602 ",
    mixBlendMode: "difference ",
  };
  const negDivInner = {
    zIndex: 10,
    backgroundColor: "#FF8E00 ",
    mixBlendMode: "difference ",
  };

  return (
    <div className="relative h-[450vh] lg:h-[350vh] overflow-hidden">
      <div className="z-20 flex flex-col h-[auto] bg-dark justify-center items-center ">
        <BackgroundShape1 />
      </div>
      <div className="z-10 flex flex-col h-[auto] bg-dark justify-center items-center ">
        <BackgroundShape2 />
      </div>
      <div className="z-0 flex flex-col h-[auto] bg-dark justify-center items-center ">
        <BackgroundShape3 />
      </div>
      <div className="absolute top-0 left-0 w-[100vw]">
        <nav className="w-[calc(100%-7px)] h-20 pl-10 flex items-center justify-end gap-10 pr-10"></nav>

        <div className="flex  lg:h-[calc(100vh-80px)] h-[calc(200vh-80px)] flex-col w-screen justify-center">
          <div className="pl-20 pt-20">
            <h1
              style={negText}
              className="text-8xl font-bold text-transparent bg-blend- font-PTMono"
            >
              TransMart
            </h1>
            <Link href={"/user/book-transport"}>
              <button
                style={negDiv}
                className="w-[auto] py-2 px-5 font-semibold hover:bg-primary/90 transition-all text-md rounded-full "
              >
                <span
                  style={{ color: "white !important", mixBlendMode: "normal" }}
                  className="flex gap-2 items-center justify-center"
                >
                  Book a Ride
                  {/* //on hover translate this to right */}
                  <ArrowRight className="transform transition-all duration-300 hover:translate-x-2" />
                </span>
              </button>
            </Link>
          </div>

          <div className="p-16 flex flex-col lg:flex-row items-center lg:justify-end gap-16">
            <div
              style={negDivInner}
              className="rounded-xl z-10 lg:w-[30%] max-w-[250px] h-[250px] bg-gray-400 flex flex-col justify-evenly items-center"
            >
              <Car className="h-20 w-20" />
              <h3 className="w-[70%] text-center text-lg">
                Book a ride from anywhere, anytime
              </h3>
            </div>
            <div
              style={negDivInner}
              className="rounded-xl z-10 lg:w-[30%] max-w-[250px] h-[250px] flex flex-col justify-evenly items-center"
            >
              <DollarSign className="h-20 w-20" />
              <h3 className="w-[70%] text-center text-lg">
                The cheapest prices you can ever get
              </h3>
            </div>
            <div
              style={negDivInner}
              className="rounded-xl z-10 lg:w-[30%] max-w-[250px] h-[250px] flex flex-col justify-evenly items-center"
            >
              <MapPinned className="h-20 w-20" />
              <h3 className="w-[70%] text-center text-lg">
                UserBase analysis for creators
              </h3>
            </div>
          </div>
        </div>

        <div className="flex justify-evenly pt-[17vh] h-[80vh]  min-w-[calc(100vw - 7px)] lg:flex-row flex-col items-center">
          <div className="lg:w-[500px] w-[80%] h-[400px] bg-gray-400 rounded-xl"></div>
          <div className="lg:w-[500px] w-[80%] h-[400px] text-center lg:text-end  lg:text-5xl text-3xl font-bold flex items-center rounded-xl">
            Stream anytime or watch your favorite creators live
          </div>
        </div>

        <div className="flex justify-evenly pt-[17vh] h-[80vh]  min-w-[calc(100vw - 7px)] lg:flex-row flex-col-reverse items-center">
          <div className="lg:w-[500px] w-[80%] h-[400px] text-center lg:text-start  lg:text-5xl text-3xl font-bold flex items-center rounded-xl">
            Invest in creators to earn when they do well
          </div>
          <div className="lg:w-[500px] w-[80%] h-[400px] bg-gray-400 rounded-xl"></div>
        </div>

        <div className="flex justify-evenly pt-[17vh] h-[80vh]  min-w-[calc(100vw - 7px)] lg:flex-row flex-col items-center">
          <div className="lg:w-[500px] w-[80%] h-[400px] bg-gray-400 rounded-xl"></div>
          <div className="lg:w-[500px] w-[80%] h-[400px] text-center lg:text-end  lg:text-5xl text-3xl font-bold flex items-center rounded-xl">
            Creators have access to specialized UserBase analysis
          </div>
        </div>
      </div>
    </div>
  );
}
