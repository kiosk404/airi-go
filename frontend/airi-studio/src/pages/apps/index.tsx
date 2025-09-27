import {LightRays} from "@/components/base/bit";

const AppsPage: React.FC = () => {
    return (
        <div style={{ width: '100%', height: '100vh', flex: 1, position: 'relative' }}>
            <LightRays
                raysOrigin="top-center"
                raysColor="#00ffff"
                raysSpeed={1.5}
                lightSpread={0.8}
                rayLength={1.2}
                followMouse={true}
                mouseInfluence={0.1}
                noiseAmount={0.1}
                distortion={0.05}
                className="custom-rays"
            />
        </div>
    );
};

export default AppsPage;