set -ex

export PATH="$HOME/compiled/bin":$PATH
export PKG_CONFIG_PATH=$HOME/compiled/lib/pkgconfig

if [ ! -e "$HOME/nasm/nasm" ]; then
  # sudo apt-get -y install asciidoc xmlto # this fails :(
  git clone -b nasm-2.14.02 https://repo.or.cz/nasm.git "$HOME/nasm"
  cd "$HOME/nasm"
  ./autogen.sh
  ./configure --prefix="$HOME/compiled"
  make
  make install || echo "Installing docs fails but should be OK otherwise"
fi

if [ ! -e "$HOME/x264/x264" ]; then
  git clone http://git.videolan.org/git/x264.git "$HOME/x264"
  cd "$HOME/x264"
  # git master as of this writing
  git checkout 545de2ffec6ae9a80738de1b2c8cf820249a2530
  ./configure --prefix="$HOME/compiled" --enable-pic --enable-static
  make
  make install-lib-static
fi

if [ ! -e "$HOME/ffmpeg/libavcodec/libavcodec.a" ]; then
  git clone https://git.ffmpeg.org/ffmpeg.git "$HOME/ffmpeg" || echo "FFmpeg dir already exists"
  cd "$HOME/ffmpeg"
  git checkout 3ea705767720033754e8d85566460390191ae27d
  ./configure --prefix="$HOME/compiled" --enable-libx264 --enable-gnutls --enable-gpl --enable-libass --enable-static --enable-libmp3lame --enable-cuda --enable-cuvid --enable-nvenc --enable-decoder=h264_cuvid 
  make
  make install
fi
