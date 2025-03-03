import sys
import cv2

def process_video(input_path, output_path):
    cap = cv2.VideoCapture(input_path)
    if not cap.isOpened():
        print("Error: Cannot open video file")
        return

    # Lấy thông tin video
    fourcc = cv2.VideoWriter_fourcc(*'mp4v')
    fps = int(cap.get(cv2.CAP_PROP_FPS))
    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))

    out = cv2.VideoWriter(output_path, fourcc, fps, (width, height))

    frame_index = 0
    while True:
        ret, frame = cap.read()
        if not ret:
            break

        # Đánh số thứ tự lên từng frame
        cv2.putText(frame, f"Frame: {frame_index}", (50, 50), cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 255, 0), 2)
        frame_index += 1
        out.write(frame)

    cap.release()
    out.release()
    print("Processing completed.")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python process_video.py <input_path> <output_path>")
        sys.exit(1)

    input_video = sys.argv[1]
    output_video = sys.argv[2]
    process_video(input_video, output_video)
